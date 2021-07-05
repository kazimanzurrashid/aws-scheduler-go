import { useState } from 'react';
import PropTypes from 'prop-types';

import CopyIcon from '@material-ui/icons/FilterNone';
import IconButton from '@material-ui/core/IconButton';
import Tooltip from '@material-ui/core/Tooltip';

const CopyToClipboardButton = ({ value }) => {
  const [show, setShow] = useState(false);

  const copy = () => {
    const temp = document.createElement('textarea');
    temp.style.position = 'absolute';
    temp.style.left = '-100%';
    temp.textContent = value;
    document.body.appendChild(temp);
    temp.select();
    document.execCommand('copy');
    document.body.removeChild(temp);
  };

  const handleClose = () => setShow(false);

  const handleClick = () => {
    setShow(true);
    copy();
  };

  return (
    <Tooltip
      onClose={handleClose}
      open={show}
      disableFocusListener
      disableTouchListener
      title="Copied"
    >
      <IconButton onClick={handleClick} onMouseOut={handleClose}>
        <CopyIcon fontSize="small"/>
      </IconButton>
    </Tooltip>
  );
};

CopyToClipboardButton.propTypes = {
  value: PropTypes.string.isRequired
}

export default CopyToClipboardButton;
