import { useState } from 'react';

import ClickAwayListener from '@material-ui/core/ClickAwayListener';
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
    <ClickAwayListener onClickAway={() => setShow(false)}>
      <Tooltip
        PopperProps={{
          disablePortal: true,
        }}
        onClose={() => setShow(false)}
        open={show}
        disableFocusListener
        disableHoverListener
        disableTouchListener
        title="Copied"
      >
        <IconButton onClick={handleClick}>
          <CopyIcon fontSize="small" />
        </IconButton>
      </Tooltip>
    </ClickAwayListener>
  );
};

export default CopyToClipboardButton;
