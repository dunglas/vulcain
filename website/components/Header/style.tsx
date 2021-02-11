import { makeStyles, Theme } from '@material-ui/core/styles';

const useStyles = makeStyles<Theme>((theme) => ({
  appBar: {
    display: 'flex',
    flexDirection: 'row',
    zIndex: theme.zIndex.drawer + 1,
  },
  mercure: {
    zIndex: theme.zIndex.drawer + 1,
    backgroundColor: '#495466',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    color: '#fff',
    fontWeight: theme.typography.fontWeightBold,
    padding: theme.spacing(1),
    width: '100%',
    '& > img': {
      height: '16px',
      margin: `0 ${theme.spacing(1)}px`,
    },
  },
  logo: {
    height: '40px',
    marginRight: '20px',
  },
  logoLink: {
    height: '100%',
    borderRadius: 0,
  },
  sponsor: {
    position: 'absolute',
    height: '40px',
    top: '46px',
    left: '70px',
    clipPath: 'polygon(2% 30%, 94% 15%, 97% 73%, 5% 79%)',
  },
  toolbar: {
    width: '100%',
    height: '64px',
  },
  menuLink: {
    height: '100%',
    borderRadius: 0,
    padding: '10px 15px 0',
    '&::after': {
      content: "''",
      position: 'absolute',
      width: '0px',
      height: '4px',
      backgroundColor: theme.palette.primary.main,
      top: 'calc(50% - 10px)',
      left: '50%',
      transform: 'translate(-50%, -50%)',
      transition: theme.transitions.create('all'),
    },
    '&:hover': {
      '&::after': {
        width: '30px',
      },
    },
    '&.active': {
      color: theme.palette.primary.main,
      position: 'relative',
      '&::after': {
        width: '20px',
      },
    },
  },
  accountButton: {
    padding: 0,
    transition: theme.transitions.create('all'),
    '&:hover': {
      color: theme.palette.primary.main,
    },
  },
  accountMenuHeader: {
    padding: theme.spacing(2, 3),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    borderBottom: `1px solid ${theme.palette.grey[200]}`,
    '& > svg': {
      fontSize: '3rem',
    },
  },
  github: {
    width: '120px',
    height: '120px',
    position: 'absolute',
    right: '0',
    top: '66px',
    overflow: 'hidden',
    clipPath: 'polygon(0% 0, 100% 0, 100% 100%)',
  },
  octoArm: {
    transformOrigin: '130px 106px',
  },
  ribbon: {
    width: '300px',
    background: theme.palette.secondary.main,
    position: 'absolute',
    textAlign: 'center',
    lineHeight: '25px',
    top: '40px',
    right: '-105px',
    transform: 'rotate(45deg)',
    zIndex: 0,
    textTransform: 'uppercase',
    color: `${theme.palette.getContrastText(theme.palette.secondary.main)}!important`,
  },
  octoSvg: {
    zIndex: 0,
    position: 'absolute',
    top: '0',
    right: '0',
    border: '0',
  },
  userMenuLink: {
    '& > a': {
      color: theme.palette.text.primary,
    },
    '& > a:hover': {
      textDecoration: 'none',
    },
  },
}));

export default useStyles;
