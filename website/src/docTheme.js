import { createTheme } from '@material-ui/core/styles';
import { red } from '@material-ui/core/colors';

const PRIMARY = '#f5731b';
const SECONDARY = '#047da7';

const initialTheme = createTheme();

// Fix anchors (because of the fixed header).
const anchorFixer = {
  '&::before': {
    content: "''",
    display: 'block',
    paddingTop: '74px',
    marginTop: '-64px',
  },
};

// Create a theme instance.
const theme = createTheme({
  palette: {
    primary: {
      main: PRIMARY,
      dark: '#e05512',
      light: '#ff9800',
    },
    secondary: {
      main: SECONDARY,
    },
    error: {
      main: red.A400,
    },
    background: {
      default: '#fff',
    },
  },
  fontFamily: "'Roboto', sans-serif",
  typography: {
    h1: {
      fontSize: '2rem',
      fontWeight: '800',
      fontFamily: "'Montserrat', sans-serif",
      borderBottom: `3px solid ${PRIMARY}`,
      color: PRIMARY,
      paddingBottom: initialTheme.spacing(1),
      marginBottom: initialTheme.spacing(1),
      '&:not(:first-of-type)': {
        marginTop: initialTheme.spacing(8),
      },
    },
    h2: {
      marginTop: initialTheme.spacing(4),
      marginBottom: initialTheme.spacing(2),
      fontSize: '1.5rem',
      fontWeight: '600',
      fontFamily: "'Montserrat', sans-serif",
      ...anchorFixer,
    },
    h3: {
      fontSize: '1.1rem',
      fontWeight: '400',
      textTransform: 'uppercase',
      fontFamily: "'Montserrat', sans-serif",
      marginTop: initialTheme.spacing(3),
      marginBottom: initialTheme.spacing(1),
      ...anchorFixer,
    },
    h4: {
      fontSize: '1rem',
      fontWeight: '400',
      fontFamily: "'Montserrat', sans-serif",
      marginTop: initialTheme.spacing(3),
      marginBottom: initialTheme.spacing(1),
      ...anchorFixer,
    },
    h5: {
      marginTop: initialTheme.spacing(2),
      marginBottom: initialTheme.spacing(1),
      fontWeight: '200',
      fontFamily: "'Montserrat', sans-serif",
      fontSize: '0.9rem',
      ...anchorFixer,
    },
    h6: {
      fontSize: '0.9rem',
      fontWeight: '800',
      ...anchorFixer,
    },
    overline: {
      lineHeight: 1.5,
    },
    body2: {
      lineHeight: 1.5,
    },
  },
  overrides: {
    MuiButton: {
      root: {
        borderRadius: '40px',
      },
    },
  },
});

export default theme;
