import React, { useEffect } from "react";
import { useTranslation } from "react-i18next";
import MainPage from "pages/MainPage.jsx";
import ExitPage from "components/Auth/ExitPage";
import SettingsPage from "pages/SettingsPage";
import AuthGuard from "components/Auth/AuthGuard";





function MainWithTitle(props) {
  const { t, i18n } = useTranslation();
  useEffect(() => {
    document.title = process.env.REACT_APP_NAME + " " + process.env.REACT_APP_VERSION;
  }, [t, i18n.language]);
  return <MainPage {...props} />;
}




export const staticRoutes = [
  { path: "/", element: <MainWithTitle /> },
  { path: "/profile", element: <AuthGuard><SettingsPage /></AuthGuard>},
  { path: "/exit", element: <ExitPage /> },
];



