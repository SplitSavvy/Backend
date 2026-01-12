import db from "../db/client.js";
import bcrypt from "bcrypt";
import "dotenv/config";
import jwt from "jsonwebtoken";
import { promisify } from "util";
import { passwordReset } from "../emails/passwordReset.js";

const saltRounds = 12;

const compareAsync = promisify(bcrypt.compare);
const hashAsync = promisify(bcrypt.hash);

//controller funciton for email verification
export const loginEmailVerify = async (req, res) => {
  const email = req.query.email;
  try {
    //checking if user exists or not
    const checkEmailPresence = await db.query(
      "SELECT * FROM users WHERE email = ($1)",
      [email]
    );
    if (checkEmailPresence.rows.length > 0) {
      //if exists moves to further step i.e. password field visibilty
      return res.status(201).json({ message: "User exist" });
    } else {
      // else shows the error page
      return res.status(404).json({ message: "User not found" });
    }
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};

//controller function for user verification
export const loginVerify = async (req, res) => {
  const { email, password } = req.body;
  try {
    const response = await db.query("SELECT * FROM users WHERE email = ($1)", [
      email,
    ]);
    if (response.rows.length === 0) {
      return res.status(404).json({ message: "User not found" });
    }
    const storedPassword = response.rows[0].hashed_password;
    const match = await compareAsync(password, storedPassword);
    if (!match) {
      return res.status(401).json({ message: "Incorrect current password" });
    }
    const token = jwt.sign({ email }, process.env.JWT_SECRET, {
      expiresIn: "30d",
    });
    return res.status(200).json({ message: "User verified", token });
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};

export const tokenVerify = async (req, res) => {
  const authHeader = req.headers["authorization"];
  const token = authHeader?.split(" ")[1];
  if (!token) {
    return res.status(400).json({ result: "Token not found" });
  }
  try {
    jwt.verify(token, process.env.JWT_SECRET);
    return res.status(200).json({ message: "User verified" });
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};

export const userDetails = async (req, res) => {
  const authHeader = req.headers["authorization"];
  const token = authHeader?.split(" ")[1];
  try {
    const decoded = jwt.verify(token, process.env.JWT_SECRET);
    const result = await db.query("SELECT * FROM users WHERE email = ($1)", [
      decoded.email,
    ]);
    const userInfo = {
      firstName: result.rows[0].first_name,
      lastName: result.rows[0].last_name,
      email: result.rows[0].email,
    };
    return res.status(200).json({ message: "User details", userInfo });
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};

export const resetPassword = async (req, res) => {
  const { email } = req.body;
  try {
    const data = await db.query("SELECT * FROM users WHERE email = ($1)", [
      email,
    ]);
    if (data.rows.length === 0) {
      return res.status(404).json({ message: "User not found." });
    }
    const user = data.rows[0];
    const token = jwt.sign(
      {
        method: "reset",
        firstName: user.first_name,
        email: user.email,
      },
      process.env.JWT_SECRET,
      { expiresIn: "15m" }
    );
    const hashed = await hashAsync(token, saltRounds);
    await db.query(
      "UPDATE users SET hashed_token = ($1), is_token_valid = ($2) WHERE email = ($3)",
      [hashed, true, email]
    );
    await passwordReset(
      email,
      user.first_name,
      `${process.env.FRONTEND_URL}/setpassword?token=${token}`
    );
    return res.status(200).json({ message: "Reset link sent" });
  } catch (err) {
    res.status(500).json({ message: "Server error" });
  }
};
