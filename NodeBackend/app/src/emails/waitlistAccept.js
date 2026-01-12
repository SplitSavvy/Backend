import resend from "./resendConfig.js";

export async function waitlistAccept(recipient, firstName) {
  try {
    const response = await resend.emails.send({
      from: "SplitSavvy <splitsavvy@savry.in>",
      to: recipient,
      subject: `✨ You’re In, ${firstName}! Your SplitSavvy adventure is about to begin 🚪🎉`,
      html: `<div style="font-family: 'Segoe UI', sans-serif; color: #333; padding: 20px; max-width: 600px; margin: auto;">
  <h2 style="color: #4CAF50;">Hey ${firstName}, guess what? You’re officially in! 🥳</h2>

  <p>Your waitlist request has been approved and we’re super excited to have you on board!</p>

  <p>In the next email, you’ll receive a special login link — that’s your golden ticket into the world of hassle-free shared living with SplitSavvy. 🧾🚪</p>

  <p>Here’s what’s waiting for you inside:</p>
  <ul>
    <li>💸 Expense splitting without the spreadsheet drama</li>
    <li>📅 Group bill tracking that won’t fry your brain</li>
    <li>📊 A dashboard so clean, Marie Kondo might cry</li>
  </ul>

  <p>Keep an eye on your inbox (and your spam folder — just in case it gets sneaky 👀).</p>

  <p>Need anything or just want to say hi? Hit us up at <a href="mailto:hello@savry.in">hello@savry.in</a> — we’re always around.</p>

  <p style="margin-top: 30px;">Let’s make shared living simple again,  
  <br><strong>— Team SplitSavvy 💚</strong></p>

  <p style="font-size: 12px; color: #888; margin-top: 20px;">
    This is an auto-generated email. Please don’t reply here — we’re likely busy arguing over who left dishes in the sink.
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
