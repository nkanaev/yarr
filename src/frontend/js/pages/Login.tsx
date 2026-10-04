import { useSignal } from "@preact/signals";
import { useI18n } from "../i18n";
import icons from "../icons";

interface LoginProps {
  onLogin: () => void;
}

export default function Login({ onLogin }: LoginProps) {
  const { t } = useI18n();
  const hasError = useSignal(false);

  const handleSubmit = (event: Event) => {
    event.preventDefault();
    const data = new FormData(event.target as HTMLFormElement);
    fetch("./login", { method: "POST", body: data }).then(res => {
      if (res.ok) {
        onLogin();
      } else {
        hasError.value = true;
      }
    });
  };

  return (
    <div className="mx-auto my-2 p-3" style={{ maxWidth: "20rem" }}>
      <form onSubmit={handleSubmit} className="d-flex flex-column">
        <div
          className="login-logo my-5 d-flex justify-content-center"
          dangerouslySetInnerHTML={{ __html: icons.anchor }}
        />
        <label htmlFor="username" className="mb-2">
          {t("username")}
        </label>
        <input
          name="username"
          className="c-input"
          id="username"
          autoComplete="off"
          required
          autoFocus
        />
        <label htmlFor="password" className="mb-2 mt-3">
          {t("password")}
        </label>
        <input name="password" className="c-input" id="password" type="password" required />
        <button className="c-button mt-3" type="submit">
          {t("login")}
        </button>
        {hasError.value && (
          <div className="fixed-top p-2 text-center bg-danger text-white">{t("login_error")}</div>
        )}
      </form>
    </div>
  );
}
