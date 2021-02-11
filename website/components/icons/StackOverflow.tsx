import React from 'react';
import { SvgIcon } from '@material-ui/core';

const SlackIcon: React.ComponentType = (props) => (
  <SvgIcon {...props}>
    <path
      d="M16.7,14.6l-9.3-2L7,14.5l9.3,1.9L16.7,14.6z M19.1,10.5l-7.3-6.1l-1.2,1.5l7.3,6.1L19.1,10.5z M17.6,12.4L9,8.4l-0.8,1.7
	l8.6,4L17.6,12.4z M15.3,1.4l-1.5,1.1l5.6,7.6L21,9L15.3,1.4z M16.3,16.9H6.8v1.9h9.5V16.9z M18.2,20.7H4.9V15H3v7.6h17V15h-1.9
	V20.7z"
    />
  </SvgIcon>
);

export default SlackIcon;
