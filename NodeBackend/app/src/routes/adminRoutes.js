import express from "express";
import {
  adminIdentiyVerify,
  tokenVerify,
  passwordUpdate,
  adminAdd,
  fetchAdminList,
  superUserPermissions,
  isSuperUser,
} from "../controllers/adminController.js";

const router = express.Router();

router.post("/verify", adminIdentiyVerify); // route to handle admin verification
router.get("/tokenverify", tokenVerify);
router.post("/passwordupdate", passwordUpdate);
router.post("/add", adminAdd); // route to add new admin details
router.post("/list", fetchAdminList);
router.post("/superuserpermissions", superUserPermissions);
router.post("/issuperuser", isSuperUser);

export default router;
