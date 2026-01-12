import express from "express";
import {
  loginEmailVerify,
  loginVerify,
  tokenVerify,
  userDetails,
  resetPassword,
} from "../controllers/authController.js";

const router = express.Router();

router.get("/check", loginEmailVerify); //route handler for email address verification
router.post("/login", loginVerify); //route handler to verify identity
router.get("/tokenverify", tokenVerify);
router.get("/userdetails", userDetails);
router.post("/reset", resetPassword);

export default router;
