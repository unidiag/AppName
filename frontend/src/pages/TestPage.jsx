import { Box, Container } from "@mui/material";
import IconButton from 'components/@extended/IconButton';
import TitleBlock from "components/TitleBlock";
import React, { useState } from "react";
import BugReportIcon from '@mui/icons-material/BugReport';
import { MinusOutlined, PlusOutlined } from "@ant-design/icons";


export default function TestPage(){

    const [count, setCount] = useState(0)

    return(
        <Container maxWidth={"xl"}>

            <TitleBlock>
                <BugReportIcon /> Test page
            </TitleBlock>

            <Box p={1}>
                this is test coutner: <b>{count}</b>
                <IconButton size={"small"} shape="rounded" variant="contained" onClick={() => setCount(count-1)}><MinusOutlined /></IconButton>
                <IconButton size={"small"} shape="rounded" variant="contained" onClick={() => setCount(count+1)}><PlusOutlined /></IconButton>
            </Box>
            


        </Container>
    )
}