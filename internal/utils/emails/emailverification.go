package emails

import "fmt"

func VerifyEmailTemplate(firstName, otp string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Welcome to HelpMe</title>
  <style>
    body {
      margin: 0;
      padding: 0;
      background-color: #f5f7fa;
      font-family: Arial, sans-serif;
      color: #333333;
    }
    .container {
      max-width: 600px;
      margin: 40px auto;
      background-color: #ffffff;
      padding: 40px 30px;
      border-radius: 8px;
      box-shadow: 0 2px 8px rgba(0,0,0,0.05);
    }
    h1 {
      color: #2c3e50;
      font-size: 24px;
      margin-bottom: 20px;
    }
    p {
      font-size: 16px;
      line-height: 1.6;
      margin-bottom: 20px;
    }
    .button {
      display: inline-block;
      padding: 14px 24px;
      background-color: #4a90e2;
      color: #ffffff;
      text-decoration: none;
      border-radius: 5px;
      font-weight: bold;
      font-size: 16px;
    }
    .footer {
      margin-top: 40px;
      font-size: 14px;
      color: #888888;
      text-align: center;
    }
    @media only screen and (max-width: 600px) {
      .container {
        padding: 30px 20px;
      }
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>Hi %s,</h1>
    <p>Welcome — we’re excited to have you on board.</p>
    <p>To activate your account, use the OTP code below:</p>
    <p style="text-align: center;">
      <button class="button">%s</button>
    </p>
    <p>This helps us keep your account secure and ensures you receive important updates.</p>
    <p>If you didn’t create this account, you can safely ignore this email — no further action is needed.</p>
    <p>See you soon,<br><strong>The Team</strong></p>
    <div class="footer">
      &copy; 2025 All rights reserved.
    </div>
  </div>
</body>
</html>
`, firstName, otp)
}
