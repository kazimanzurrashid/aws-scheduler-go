import { useState } from 'react';

import Tooltip from '@material-ui/core/Tooltip';
import IconButton from '@material-ui/core/IconButton';
import CopyIcon from '@material-ui/icons/FilterNone';

const CopyToClipboardButton = ({ onClick }) => {
  const [show, setShow] = useState(false);

  const handleClick = () => {
    setShow(true);
    onClick();
  };

  return (
    <Tooltip
      onClose={() => setShow(false)}
      open={show}
      disableFocusListener
      disableHoverListener
      disableTouchListener
      title="Copied"
    >
      <IconButton onClick={handleClick} onMouseOut={() => setShow(false)}>
        <CopyIcon fontSize="small" />
      </IconButton>
    </Tooltip>
  );
};

export default CopyToClipboardButton;
