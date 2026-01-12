import resend from "./resendConfig.js";

export async function firstLogin(recipient, firstName, link) {
  try {
    const response = await resend.emails.send({
      from: "SplitSavvy <splitsavvy@savry.in>",
      to: recipient,
      subject: `🔓 Your SplitSavvy spot is unlocked — set your password & jump in!`,
      html: `<div style="font-family: 'Segoe UI', sans-serif; color: #333; padding: 20px; max-width: 600px; margin: auto;">
  <h2 style="color: #4CAF50;">Hey ${firstName}, welcome to the club! 🎉</h2>

  <p>You've officially been accepted into SplitSavvy. It's time to claim your spot and set your password.</p>

  <p>🔐 This is a <strong>one-time login link</strong> that'll take you straight to your setup page:</p>

  <div style="margin: 20px 0;">
    <a href="${link}" style="background-color: #4CAF50; color: white; padding: 12px 20px; text-decoration: none; border-radius: 5px; font-weight: bold;">
      Set My Password
    </a>
  </div>

  <p><strong>This link will expire in 7 days</strong> — so don’t ghost us 👻</p>

  <p>If the button above doesn't work, you can copy and paste this link into your browser:</p>
  <p style="font-size: 14px; color: #555;">${link}</p>

  <p>If you weren’t expecting this email, you can safely ignore it. No changes will be made unless you click the link.</p>

  <p style="margin-top: 30px;">Catch you inside,  
  <br><strong>Team SplitSavvy 💚</strong></p>

  <p style="font-size: 12px; color: #888; margin-top: 20px;">
    This is an auto-generated email — replying won't do much (but we still love you).
  </p>
</div>

`,
    });
    console.log(`✅ Confirmation email sent to ${recipient}`);
    console.log("📨 Email response:", response);
  } catch (err) {
    console.error(`❌ Failed to send confirmation email to ${recipient}`);
    console.error("Error details:", err);
  }
}
