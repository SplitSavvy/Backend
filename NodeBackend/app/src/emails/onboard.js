import resend from "./resendConfig.js";

export async function onboard(recipient, firstName) {
  try {
    const response = await resend.emails.send({
      from: "SplitSavvy <splitsavvy@savry.in>",
      to: recipient,
      subject: `🚀 Welcome to SplitSavvy, ${firstName}! Let’s get you set up`,
      html: `<div style="font-family: 'Segoe UI', sans-serif; color: #333; padding: 20px; max-width: 600px; margin: auto;">
  <h2 style="color: #4CAF50;">Welcome aboard, ${firstName}! 🎉</h2>

  <p>We're stoked to have you here. You've just unlocked the smarter, sassier way to handle shared living — minus the chaos and awkward money talks.</p>

  <p>Here’s what you can dive into right now:</p>
  <ul>
    <li>🏠 Create or join your household crew</li>
    <li>🧾 Add shared expenses (no more “who owes what?” drama)</li>
    <li>📦 (Coming soon!) Track inventory so the snacks never run out</li>
    <li>🧮 Sit back while we crunch the numbers — effortlessly</li>
  </ul>

  <p>Go ahead, log in, explore the dashboard, and make yourself at home. 🛋️  
  Got questions or just wanna say hey? Ping us at <a href="mailto:hello@savry.in">hello@savry.in</a>.</p>

  <p>We’re always building new features and making things better, so stick around — the good stuff is just getting started. ✨</p>

  <p style="margin-top: 30px;">Happy splitting,  
  <br><strong>— Team SplitSavvy 💚</strong></p>

  <p style="font-size: 12px; color: #888; margin-top: 20px;">
    This is an auto-generated email. Replies won’t reach us (we're probably buried in a pile of receipts).
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
