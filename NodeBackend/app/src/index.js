import express from "express";
import cors from "cors";
import waitlistRoutes from "./routes/waitlistRoutes.js";
import authRoutes from "./routes/authRoutes.js";
import adminRoutes from "./routes/adminRoutes.js";

//defining express app and port
const app = express();
const port = 3000;

//middlewares cors to enable cross origin resource sharing express.json to parse req
app.use(cors());
app.use(express.json());

//route handlers
app.use("/waitlist", waitlistRoutes);
app.use("/auth", authRoutes);
app.use("/admin", adminRoutes);

//app start and confirmation log outpt
app.listen(port, () => {
  console.log(`Server running on port ${port}`);
});
