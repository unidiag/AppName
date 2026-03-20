import React, { useEffect, useState } from "react";
import AuthGuard from "components/Auth/AuthGuard";
import { sendDataToServer } from "utils/functions";
import { Box } from "@mui/material";



export default function MainPage(){

    const [time, setTime] = useState("")




    useEffect(() => {
        const loadTime = () => {
            sendDataToServer({ op: "getTime" }).then((res) => {
                if (res.status === "OK") {
                    setTime(res.time);
                }
            });
        };

        loadTime();
        const interval = setInterval(loadTime, 1000);

        return () => clearInterval(interval);
    }, []);





    return (
        <AuthGuard>
            
            <Box
                sx={{
                    height: "100vh",
                    display: "flex",
                    justifyContent: "center", // horizontal center
                    alignItems: "center",     // vertical center
                    fontSize: "32px"
                }}
            >
                {time}
            </Box>

        </AuthGuard>
    )

}