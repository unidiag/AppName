import { Box, Button } from "@mui/material";
import TitleBlock from "components/TitleBlock";
import React, { useState } from "react";
import BugReportIcon from '@mui/icons-material/BugReport';


export default function TestPage(){

    const [count, setCount] = useState(0)

    return(
        <Box px={2}>
            <TitleBlock>
                <BugReportIcon /> Test page
            </TitleBlock>

            <Box p={1}>
                this is test coutner: <b>{count}</b>
                <Button onClick={() => setCount(count+1)}>+1</Button>
                <Button onClick={() => setCount(count-1)}>-1</Button>
            </Box>

        </Box>
    )
}