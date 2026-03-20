import AuthGuard from "components/Auth/AuthGuard";
import SettingsBlock from "components/Settings/SettingsBlock";
import UsersBlock from "components/Settings/UsersBlock";
import { useAuth } from "utils/useAuth";



export default function SettingsPage() {

  const {user} = useAuth()
  const readonly = user?.status !== 9

  return (
    <AuthGuard>
      {!readonly && <SettingsBlock />}
      <UsersBlock readonly={readonly}  />
    </AuthGuard>
  );
}