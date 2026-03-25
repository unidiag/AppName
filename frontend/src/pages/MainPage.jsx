import CottageOutlinedIcon from '@mui/icons-material/CottageOutlined';
import { Container, FormControl, FormControlLabel, InputLabel, MenuItem, Select, Switch } from "@mui/material";
import TitleBlock from "components/TitleBlock";
import React, { useState } from "react";
import { useTranslation } from 'react-i18next';

export default function MainPage() {

  const {t,} = useTranslation()
  const [age, setAge] = useState(20)

  return (
      <Container maxWidth="xl">
        <TitleBlock
          t1={<FormControlLabel control={<Switch defaultChecked color="success" size="small" />} label="t1" />}
          t2={
            <FormControl>
              <InputLabel id="demo-simple-select-label">Select t2</InputLabel>
              <Select labelId="demo-simple-select-label" id="demo-simple-select" value={age} placeholder="t2" onChange={e => setAge(e.target.value)}>
                <MenuItem value={10}>Ten</MenuItem>
                <MenuItem value={20}>Twenty</MenuItem>
                <MenuItem value={30}>Thirty</MenuItem>
              </Select>
            </FormControl>
          }
          t3={<span>t3</span>}          
        >
          <CottageOutlinedIcon /> {t("home")}
        </TitleBlock>
      </Container>
  );
}