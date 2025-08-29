package templates

import (
	"fmt"
	"html"
)

func RecoverPasswordHTML(appName, username, resetLink string) string {
	escApp := html.EscapeString(appName)
	escUser := html.EscapeString(username)
	escLink := html.EscapeString(resetLink)

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="pt-BR">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Recuperação de senha - %s</title>
  <style>
    body { background:#f6f9fc; margin:0; padding:0; font-family: Arial, Helvetica, sans-serif; color:#0f172a; }
    .container { max-width:600px; margin:0 auto; padding:24px; }
    .card { background:#ffffff; border-radius:12px; padding:24px; box-shadow:0 1px 4px rgba(0,0,0,0.06); }
    h1 { font-size:20px; margin:0 0 12px; }
    p { line-height:1.6; margin:0 0 12px; }
    .btn { display:inline-block; padding:12px 16px; background:#2563eb; color:#ffffff !important; border-radius:8px; text-decoration:none; font-weight:bold; }
    .muted { color:#64748b; font-size:12px; }
    .footer { text-align:center; margin-top:16px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="card">
      <h1>Recuperação de senha</h1>
      <p>Olá, %s.</p>
      <p>Recebemos uma solicitação para redefinir a sua senha no %s.</p>
      <p>Para criar uma nova senha, clique no botão abaixo:</p>
      <p>
        <a class="btn" href="%s" target="_blank" rel="noopener">Redefinir senha</a>
      </p>
      <p class="muted">Se você não solicitou a alteração, ignore este e-mail. O link expira automaticamente após algum tempo por segurança.</p>
    </div>
    <div class="footer">
      <p class="muted">%s</p>
    </div>
  </div>
</body>
</html>`, escApp, escUser, escApp, escLink, escApp)
}
