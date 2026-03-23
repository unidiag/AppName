import { Container } from "@mui/material";
import SettingsBlock from "components/Settings/SettingsBlock";
import UsersBlock from "components/Settings/UsersBlock";
import { useAuth } from "utils/useAuth";



export default function SettingsPage() {

  const {user} = useAuth()
  const readonly = user?.status !== 9

  return (
    <Container maxWidth="xl">
        {!readonly && <SettingsBlock />}
        <UsersBlock readonly={readonly}  />
    </Container>
  );
}