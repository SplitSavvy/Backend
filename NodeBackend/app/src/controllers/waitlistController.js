import db from "../db/client.js";
import jwt from "jsonwebtoken";
import bcrypt from "bcrypt";
import "dotenv/config";
import { promisify } from "util";
import { waitlistConfirmation } from "../emails/waitlistConfirmation.js";
import { waitlistReject } from "../emails/waitlistReject.js";
import { firstLogin } from "../emails/firstLogin.js";
import { waitlistAccept } from "../emails/waitlistAccept.js";
import { onboard } from "../emails/onboard.js";

const saltRounds = 12;

const compareAsync = promisify(bcrypt.compare);
const hashAsync = promisify(bcrypt.hash);

//controller function to add users to waitlist
export const addToWaitlist = async (req, res) => {
  const { firstName, lastName, email } = req.body;
  try {
    //checking if user already in waitlist{edgecase}
    const checkWaitlist = await db.query(
      "SELECT * FROM waitlists WHERE email = ($1)",
      [email]
    );
    if (checkWaitlist.rows.length > 0) {
      return res.status(409).json({ message: "Already waitlisted" }); // shows already in waitlist{edgecase}
    }
    // checking if user is already a member{edgecase}
    const checkUserExist = await db.query(
      "SELECT * FROM users WHERE email = ($1)",
      [email]
    );
    if (checkUserExist.rows.length > 0) {
      return res.status(409).json({ message: "User exists" }); //shows already a member{edgecase}
    }
    await db.query(
      "INSERT INTO waitlists (first_name, last_name, email) VALUES ($1, $2, $3)", //add to waitlist
      [firstName, lastName, email]
    );
    await waitlistConfirmation(email, firstName);
    return res.status(201).json({ message: "Waitlisted" });
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};

//controller function to check waitlist status
export const checkInWaitlist = async (req, res) => {
  const email = req.query.email;
  try {
    const checkUserExist = await db.query(
      //checking if user already a memeber{edgecase}
      "SELECT * FROM users WHERE email = ($1)",
      [email]
    );
    const checkStatus = await db.query(
      //checking the user status
      "SELECT * FROM waitlists WHERE email = ($1)",
      [email]
    );
    if (checkStatus.rows.length === 0) {
      return res.status(404).json({ message: "User not on waitlist" }); //error when user checking status without joining the waitlist {edgecase}
    }
    if (checkUserExist.rows.length > 0) {
      return res.status(409).json({ message: "User exists" }); //returns as user already a memeber {edgecase}
    }
    //returns the status
    const row = checkStatus.rows[0];
    const fullName = `${row.first_name} ${row.last_name}`;
    console.log(row.status);
    console.log(fullName);
    return res.status(200).json({
      status: row.status,
      name: fullName,
    });
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};

export const fetchWaitlistData = async (req, res) => {
  try {
    const response = await db.query(
      "SELECT * FROM waitlists ORDER BY CASE status WHEN 'pending' THEN 1 WHEN 'approved' THEN 2 WHEN 'cancelled' THEN 3 END"
    );
    const result = response.rows;
    res.status(200).json({ message: "Waitlist", result });
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};

export const approveWaitlist = async (req, res) => {
  const { email } = req.body;
  try {
    const response = await db.query(
      "UPDATE waitlists SET status = ($1) WHERE email = ($2) RETURNING *",
      ["approved", email]
    );
    const userDetails = {
      firstName: response.rows[0].first_name,
      lastName: response.rows[0].last_name,
      email: response.rows[0].email,
    };
    await waitlistAccept(userDetails.email, userDetails.firstName);
    const token = jwt.sign(
      {
        method: "create",
        firstName: userDetails.firstName,
        lastName: userDetails.lastName,
        email: userDetails.email,
      },
      process.env.JWT_SECRET,
      {
        expiresIn: "7d",
      }
    );
    const hashed = await hashAsync(token, saltRounds);
    await db.query(
      "UPDATE waitlists SET hashed_token = ($1), is_token_valid = ($2) WHERE email = ($3)",
      [hashed, true, userDetails.email]
    );
    await firstLogin(
      userDetails.email,
      userDetails.firstName,
      `${process.env.FRONTEND_URL}/setpassword?token=${token}`
    );
    return res.status(200).json({ message: "Successfully accepted" });
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};

export const rejectWaitlist = async (req, res) => {
  const { email } = req.body;
  try {
    const response = await db.query(
      "UPDATE waitlists SET status = ($1) WHERE email = ($2) RETURNING first_name",
      ["cancelled", email]
    );
    const firstName = response.rows[0].first_name;
    await waitlistReject(email, firstName);
    return res.status(200).json({ result: "Successfully rejected" });
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};

export const tokenValidation = async (req, res) => {
  const authHeader = req.headers["authorization"];
  const token = authHeader?.split(" ")[1];
  let decoded;
  try {
    decoded = jwt.verify(token, process.env.JWT_SECRET);
  } catch (err) {
    return res.status(401).json({ message: "Invalid or expired token" });
  }
  try {
    const response = await db.query(
      "SELECT * FROM waitlists WHERE email = ($1)",
      [decoded.email]
    );
    const match = await compareAsync(token, response.rows[0].hashed_token);
    if (!match) {
      return res.status(401).json({ message: "Token mismatch" });
    }
    if (match && response.rows[0].is_token_valid) {
      const mode = decoded.method;
      return res.status(200).json({ message: "Success", mode });
    } else {
      return res
        .status(401)
        .json({ message: "Link no longer valid. Contact support for help." });
    }
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};

export const onBoardUser = async (req, res) => {
  const { token, password } = req.body;
  let decoded;
  try {
    decoded = jwt.verify(token, process.env.JWT_SECRET);
  } catch (err) {
    return res.status(401).json({ message: "Some error occured" });
  }
  try {
    const validityCheck = await db.query(
      "SELECT is_token_valid FROM waitlists WHERE email = ($1)",
      [decoded.email]
    );
    if (validityCheck.rows[0].is_token_valid) {
      const hashed = await hashAsync(password, saltRounds);
      const response = await db.query(
        "INSERT INTO users (email, first_name, last_name, hashed_password) VALUES ($1, $2, $3, $4) RETURNING *",
        [decoded.email, decoded.firstName, decoded.lastName, hashed]
      );
      if (response.rows.length === 0) {
        return res
          .status(500)
          .json({ message: "Failed to update details to users" });
      } else {
        await db.query("DELETE FROM waitlists WHERE email=($1)", [
          decoded.email,
        ]);
        await onboard(decoded.email, decoded.firstName);
      }
      return res.status(200).json({ message: "Success" });
    } else {
    }
  } catch (err) {
    return res.status(500).json({ message: "Server error" });
  }
};
