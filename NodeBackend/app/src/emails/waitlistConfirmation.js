import resend from "./resendConfig.js";

export async function waitlistConfirmation(recipient, firstName) {
  try {
    const response = await resend.emails.send({
      from: "SplitSavvy <splitsavvy@savry.in>",
      to: recipient,
      subject: `🎉 You're officially waitlisted, ${firstName}!`,
      html: `<div style="font-family: 'Segoe UI', sans-serif; color: #333; padding: 20px; max-width: 600px; margin: auto;">
  <h2 style="color: #4CAF50;">Hey ${firstName}, you made it! 🎉</h2>

  <p>Welcome aboard the SplitSavvy train — the smartest way to manage shared living without losing your sanity (or socks 🧦).</p>

  <p>You're officially on the waitlist, and we’re stoked to have you. Here’s what’s coming your way real soon:</p>
  <ul>
    <li>💸 Magic-like expense tracking</li>
    <li>🔄 Auto-split bills (no awkward convos needed)</li>
    <li>📦 Low-inventory alerts so you never run out of snacks again</li>
    <li>📊 A super clean dashboard to keep life under control</li>
  </ul>

  <p>⏳ We usually approve new users within a week — you’ll hear from us by <strong>next Monday</strong> at the latest!</p>

  <p style="font-style: italic; color: #666;">P.S. This is an auto-generated email, so replying won’t work (our inbox is probably snacking somewhere).</p>

  <p>If you need help or just want to say hi, reach us at <a href="mailto:hello@savry.in">hello@savry.in</a> — we’d love to hear from you!</p>

  <p style="margin-top: 30px;">Stay awesome,  
  <br><strong>Team SplitSavvy 💚</strong></p>

  <p style="font-size: 12px; color: #888; margin-top: 20px;">
    You received this email because you signed up on SplitSavvy. If it wasn’t you, feel free to ignore it — no hard feelings!
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
