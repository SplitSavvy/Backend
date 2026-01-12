import db from "../db/client.js";
import bcrypt from "bcrypt";
import "dotenv/config";
import jwt from "jsonwebtoken";
import { promisify } from "util";

const saltRounds = 12;

const compareAsync = promisify(bcrypt.compare);
const hashAsync = promisify(bcrypt.hash);

export const adminIdentiyVerify = async (req, res) => {
  const { username, password } = req.body;
  const lowercaseUsername = username.toLowerCase();
  try {
    const response = await db.query(
      "SELECT * FROM admins where username = ($1)",
      [lowercaseUsername]
    );
    if (response.rows.length === 0) {
      return res.status(404).json({ message: "Admin not found" });
    }
    const storedPassword = response.rows[0].hashed_password;
    bcrypt.compare(password, storedPassword, (err, result) => {
      if (err) {
        return res.status(500).json({ message: "Password comparision failed" });
      } else {
        try {
          if (result) {
            const token = jwt.sign({ username }, process.env.JWT_SECRET, {
              expiresIn: "30m",
            });
            res.status(200).json({ message: "Admin Verified", token });
          } else {
            res.status(401).json({ message: "Incorrect password" });
          }
        } catch (err) {
          res.status(500).json({ message: "Token generation failed" });
        }
      }
    });
  } catch (err) {
    res.status(500).json({ message: "Server error" });
  }
};

export const tokenVerify = async (req, res) => {
  const authHeader = req.headers["authorization"];
  const token = authHeader?.split(" ")[1];

  if (!token) return res.status(404).json({ message: "No token" });

  try {
    jwt.verify(token, process.env.JWT_SECRET);
    res.status(200).json({ message: "Verified" });
  } catch (err) {
    res.status(401).json({ message: "Invalid or expired token" });
  }
};

export const passwordUpdate = async (req, res) => {
  const { oldPassword, newPassword, token } = req.body;
  let username;
  try {
    const decoded = jwt.verify(token, process.env.JWT_SECRET);
    username = decoded.username;
  } catch (err) {
    return res.status(401).json({ message: "Invalid token" });
  }
  try {
    const response = await db.query(
      "SELECT * FROM admins where username = ($1)",
      [username]
    );
    if (response.rows.length === 0) {
      return res.status(404).json({ message: "Admin not found" });
    }
    const storedPassword = response.rows[0].hashed_password;
    const match = await compareAsync(oldPassword, storedPassword);
    if (!match) {
      return res.status(401).json({ message: "Incorrect current password" });
    }
    const hashed = await hashAsync(newPassword, saltRounds);
    await db.query(
      "UPDATE admins SET hashed_password = ($1) WHERE username = ($2)",
      [hashed, username]
    );
    return res.status(200).json({ message: "Password updated Successfully" });
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};

export const adminAdd = async (req, res) => {
  const { username, password, superUser } = req.body;
  const lowercaseUsername = username.toLowerCase();
  try {
    const usernameCheck = await db.query(
      "SELECT * FROM admins WHERE username = ($1)",
      [lowercaseUsername]
    );
    if (usernameCheck.rows.length > 0) {
      return res.status(409).json({ message: "Username not available" });
    }
    const hashed = await hashAsync(password, saltRounds);
    await db.query(
      "INSERT INTO admins (username, hashed_password, super_user)VALUES ($1, $2, $3)",
      [lowercaseUsername, hashed, superUser]
    );
    return res.status(201).json({ message: "Admin created" });
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};

export const fetchAdminList = async (req, res) => {
  try {
    const response = await db.query(
      "SELECT admin_id, username, super_user, super_admin FROM admins ORDER BY admin_id"
    );
    const list = response.rows;
    return res.status(200).json({ message: "Admin List", list });
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};

export const superUserPermissions = async (req, res) => {
  const { username, superUser } = req.body;
  try {
    await db.query(
      "UPDATE admins SET super_user = ($1) WHERE username = ($2)",
      [superUser, username]
    );
    res.status(200).json({ message: "Success" });
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};

export const isSuperUser = async (req, res) => {
  const { token } = req.body;
  const data = jwt.verify(token, process.env.JWT_SECRET);
  const username = data.username;
  try {
    const response = await db.query(
      "SELECT super_user FROM admins WHERE username = ($1)",
      [username]
    );
    const superUser = response.rows[0].super_user;
    res.status(200).json({ message: "Success", superUser });
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};
