import { MenuItem, Select, TableCell, Typography } from '@mui/material';
import { useState } from 'react';
import type { SelectChangeEvent } from '@mui/material';

interface Props<T> {
  value: T;
  list: T[];
  handleSave: (value: T) => void;
}

const CustomTableCellPulldown = <T extends string | boolean>({
  value,
  list,
  handleSave,
}: Props<T>) => {
  const [editFlag, setEditFlag] = useState<boolean>(false);

  const onChange = (e: SelectChangeEvent) => {
    switch (typeof value) {
      case 'string': {
        handleSave(e.target.value as T);
        break;
      }
      case 'boolean': {
        const bool = e.target.value === 'true' ? true : false;
        handleSave(bool as T);
        break;
      }
    }
  };

  return (
    <TableCell onClick={() => setEditFlag(true)}>
      {editFlag ? (
        <Select
          id="body"
          value={value.toString()}
          onChange={onChange}
          onBlur={() => setEditFlag(false)}
        >
          {list.map((value, index) => (
            <MenuItem
              key={index}
              value={value.toString()}
            >
              {value.toString()}
            </MenuItem>
          ))}
        </Select>
      ) : (
        <Typography>{value.toString()}</Typography>
      )}
    </TableCell>
  );
};

export default CustomTableCellPulldown;
