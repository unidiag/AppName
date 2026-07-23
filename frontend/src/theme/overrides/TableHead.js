// ==============================|| OVERRIDES - TABLE HEAD ||============================== //

export default function TableHead(theme) {
  return {
    MuiTableHead: {
      styleOverrides: {
        root: {
          backgroundColor: theme.palette.background.default,
          borderTop: `1px solid ${theme.palette.divider}`,
          borderBottom: `2px solid ${theme.palette.divider}`,

          "& .MuiTableCell-head": {
            color: theme.palette.text.primary,
            fontWeight: 600,
          },
        },
      },
    },
  };
}
