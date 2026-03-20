import PropTypes from "prop-types";
import { Link as RouterLink } from "react-router-dom";
import { useTheme } from "@mui/material/styles";
import { Box } from "@mui/material";

export default function MyLink({
  to,
  href,
  children,
  target = "_self",
  underline = false,
  sx = {},
  title,
  ...rest
}) {
  const theme = useTheme();

  const color =
    theme.palette.mode === "dark"
      ? theme.palette.primary.light
      : theme.palette.primary.main;

  const commonSx = {
    color,
    textDecoration: underline ? "underline" : "none",
    cursor: "pointer",
    "&:hover": {
      textDecoration: "underline",
    },
    ...sx,
  };

  // external
  if (href) {
    return (
      <Box
        component="a"
        href={href}
        target={target}
        rel="noopener noreferrer"
        sx={commonSx}
        title={title}
        {...rest}
      >
        {children}
      </Box>
    );
  }

  // internal
  return (
    <Box
      component={RouterLink}
      to={to}
      target={target}
      sx={commonSx}
      title={title}
      {...rest}
    >
      {children}
    </Box>
  );
}

MyLink.propTypes = {
  to: PropTypes.string,
  href: PropTypes.string,
  children: PropTypes.node,
  target: PropTypes.string,
  underline: PropTypes.bool,
  sx: PropTypes.object,
  title: PropTypes.string,
};