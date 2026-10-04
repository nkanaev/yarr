import { render } from "preact";
import { useSignal } from "@preact/signals";
import type { Lang } from "./i18n";
import { I18nProvider } from "./i18n";
import App from "./pages/App";
import Login from "./pages/Login";
import api from "./api";

function Root() {
  const authenticated = useSignal(window.app.authenticated);

  const onLogin = () => {
    api.settings.get().then(settings => {
      window.app.settings = settings;
      window.app.authenticated = true;
      authenticated.value = true;
    });
  };

  const initialLang = (window.app.settings?.language ||
    document.documentElement.lang ||
    "en") as Lang;

  return (
    <I18nProvider initialLang={initialLang}>
      {authenticated.value ? <App /> : <Login onLogin={onLogin} />}
    </I18nProvider>
  );
}

const rootEl = document.getElementById("app");
if (rootEl) {
  render(<Root />, rootEl);
}
