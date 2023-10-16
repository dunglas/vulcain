import { createTheme } from '@material-ui/core/styles';
import { red } from '@material-ui/core/colors';

const PRIMARY = '#f5731b';
const SECONDARY = '#047da7';
const TEXT_PRIMARY = 'rgba(0, 0, 0, 0.87);';

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
    textPrimary: {
      main: TEXT_PRIMARY,
    },
  },
  fontFamily: "'Roboto', sans-serif",
  typography: {
    h1: {
      fontSize: '4.5rem',
      fontWeight: '800',
      fontFamily: "'Montserrat', sans-serif",
    },
    h2: {
      fontSize: '3.2rem',
      fontWeight: '800',
      fontFamily: "'Montserrat', sans-serif",
    },
    h3: {
      fontSize: '2.5rem',
      fontWeight: '800',
      fontFamily: "'Montserrat', sans-serif",
    },
    h5: {
      fontWeight: '400',
      fontFamily: "'Montserrat', sans-serif",
      fontSize: '1.2rem',
      lineHeight: 1,
      textTransform: 'uppercase',
    },
    h6: {
      fontSize: '1.1rem',
    },
    overline: {
      lineHeight: 1.5,
    },
  },
  overrides: {
    MuiCssBaseline: {
      '@global': {
        ':target::before': {
          content: '""',
          display: 'block',
          height: '64px' /* fixed header height*/,
          margin: '-64px 0 0' /* negative fixed header height */,
        },
      },
    },
    MuiButton: {
      root: {
        borderRadius: '40px',
      },
    },
    MuiAppBar: {
      root: {
        borderBottom: `2px solid ${PRIMARY}`,
      },
      colorDefault: {
        backgroundColor: '#fff',
        color: TEXT_PRIMARY,
      },
    },
  },
});

export default theme;
