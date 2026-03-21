import React, { useState } from "react";
import AuthGuard from "components/Auth/AuthGuard";
import { sendDataToServer } from "utils/functions";
import {
  Box,
  TextField,
  Button,
  Paper
} from "@mui/material";
import MarkdownRenderer from "components/MarkdownRenederer";

export default function MainPage() {
  const [prompt, setPrompt] = useState("");
  const [answer, setAnswer] = useState("");
  const [loading, setLoading] = useState(false);

  const handleAsk = () => {
    if (!prompt.trim()) return;

    setLoading(true);
    setAnswer("");

    sendDataToServer({
      op: "askAI",
      prompt
    }).then((res) => {
      if (res.status === "OK") {
        setAnswer(res.answer || "");
      } else {
        setAnswer("Error: " + res.status);
      }
      setLoading(false);
    });
  };

  return (
    <AuthGuard>
      <Box
        sx={{
          height: "80vh",
          display: "flex",
          flexDirection: "column",
          p: 2,
          gap: 2
        }}
      >
        {/* Answer */}
        <Paper
            sx={{
                flex: 1,
                p: 2,
                overflow: "auto"
            }}
        >
            {loading ? "Thinking..." : <MarkdownRenderer content={answer} />}
        </Paper>

        {/* Input */}
        <Box sx={{ display: "flex", gap: 2 }}>
          <TextField
            fullWidth
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            placeholder="Enter your prompt..."
          />
          <Button
            variant="contained"
            onClick={handleAsk}
            disabled={loading}
          >
            Send
          </Button>
        </Box>
      </Box>
    </AuthGuard>
  );
}