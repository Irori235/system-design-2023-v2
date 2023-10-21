import { TableCell, TextField, Typography } from '@mui/material';
import { useState, useRef } from 'react';
import type { FC } from 'react';

interface Props {
  value: string;
  handleSave: (value: string) => void;
}

const CustomTableCell: FC<Props> = ({ value, handleSave }) => {
  const [editFlag, setEditFlag] = useState<boolean>(false);
  const ref = useRef<HTMLInputElement>(null);

  const onClick = () => {
    setEditFlag(true);
  };

  const onBlur = () => {
    if (ref.current?.value) {
      handleSave(ref.current.value);
    }
    setEditFlag(false);
  };

  return (
    <TableCell onClick={onClick}>
      {editFlag ? (
        <TextField
          id="body"
          inputRef={ref}
          defaultValue={value}
          onBlur={onBlur}
          autoFocus
        />
      ) : (
        <Typography>{value}</Typography>
      )}
    </TableCell>
  );
};

export default CustomTableCell;
