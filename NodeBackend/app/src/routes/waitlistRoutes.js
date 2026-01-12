import express from "express";
import {
  addToWaitlist,
  checkInWaitlist,
  fetchWaitlistData,
  approveWaitlist,
  rejectWaitlist,
  tokenValidation,
  onBoardUser,
} from "../controllers/waitlistController.js";

const router = express.Router();

router.post("/add", addToWaitlist); //route handler to add and check for edge cases in waitlist
router.get("/check", checkInWaitlist); //route handler to check the status in waitlist
router.post("/list", fetchWaitlistData);
router.post("/reject", rejectWaitlist);
router.post("/approve", approveWaitlist);
router.get("/tokenverify", tokenValidation);
router.post("/onboard", onBoardUser);

export default router;
