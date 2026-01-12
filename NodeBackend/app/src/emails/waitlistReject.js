import resend from "./resendConfig.js";

export async function waitlistReject(recipient, firstName) {
  try {
    const response = await resend.emails.send({
      from: "SplitSavvy <splitsavvy@savry.in>",
      to: recipient,
      subject: `😔 About Your SplitSavvy Invite…`,
      html: `<div style="font-family: 'Segoe UI', sans-serif; color: #333; padding: 20px; max-width: 600px; margin: auto;">
  <h2 style="color: #FF6B6B;">Hey ${firstName},</h2>

  <p>Thanks a ton for joining the SplitSavvy waitlist — we truly appreciate your interest in making shared living smarter!</p>

  <p>Unfortunately, due to a high volume of registrations this week, we weren't able to approve your invite this time around. 😞</p>

  <p>But don’t worry — this isn’t the end! You can try again after <strong>48 hours</strong>, and we’d love to see you back in the queue.</p>

  <p>If you have any questions or just want to say hi, feel free to reach us at <a href="mailto:hello@savry.in">hello@savry.in</a>.</p>

  <p style="margin-top: 30px;">All the best,  
  <br><strong>Team SplitSavvy 💚</strong></p>

  <p style="font-size: 12px; color: #888; margin-top: 20px;">
    This is an auto-generated email — no one’s actually reading replies here (we promise we’re not being rude).
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
