import resend from "./resendConfig.js";

export async function passwordReset(recipient, firstName, link) {
  try {
    const response = await resend.emails.send({
      from: "SplitSavvy <splitsavvy@savry.in>",
      to: recipient,
      subject: `🛠️ Forgot your password? Let’s fix that, ${firstName}!`,
      html: `<div style="font-family: 'Segoe UI', sans-serif; color: #333; padding: 20px; max-width: 600px; margin: auto;">
  <h2 style="color: #4CAF50;">Hey ${firstName}, no worries — it happens! 🔑</h2>

  <p>We all blank out sometimes. That’s why we’ve got your back with a secure, one-time link to reset your SplitSavvy password.</p>

  <p>Click the button below to choose a shiny new password:</p>

  <div style="margin: 20px 0;">
    <a href="${link}" style="background-color: #4CAF50; color: white; padding: 12px 20px; text-decoration: none; border-radius: 5px; font-weight: bold;">
      Reset My Password
    </a>
  </div>

  <p><strong>This link is valid for the next 15 minutes</strong> — after that, it disappears like the last cookie in the pantry 🍪</p>

  <p>If the button doesn’t work, you can copy and paste this into your browser:</p>
  <p style="font-size: 14px; color: #555;">${link}</p>

  <p>If you didn’t request a reset, feel free to ignore this email. Your account is safe and sound.</p>

  <p style="margin-top: 30px;">See you back inside,  
  <br><strong>Team SplitSavvy 💚</strong></p>

  <p style="font-size: 12px; color: #888; margin-top: 20px;">
    This is an auto-generated email — replies won’t be seen. For real help, hit us up at <a href="mailto:hello@savry.in">hello@savry.in</a>.
  </p>
{l
`,
    });
    console.log(`✅ Confirmation email sent to ${recipient}`);
    console.log("📨 Email response:", response);
  } catch (err) {
    console.error(`❌ Failed to send confirmation email to ${recipient}`);
    console.error("Error details:", err);
  }
}
