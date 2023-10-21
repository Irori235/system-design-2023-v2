import { Delete } from '@mui/icons-material';
import { Button, IconButton, TextField } from '@mui/material';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import React, { useRef, useState } from 'react';
import CustomTableCell from './CustomTableCell';
import CustomTableCellPulldown from './CustomTableCellPulldown';
import type { Task } from '@/types/tasks';
import type { FC } from 'react';

interface Props {
  tasks: Task[];
  postTask: (title: string) => Promise<void>;
  updateTask: (id: string, title: string, isDone: boolean) => Promise<void>;
  deleteTask: (id: string) => Promise<void>;
}

const TaskTable: FC<Props> = ({ tasks, postTask, updateTask, deleteTask }) => {
  const [isAddTaskFlag, setIsAddTaskFlag] = useState<boolean>(false);
  const titleRef = useRef<HTMLInputElement>(null);

  const onBlurTitle = () => {
    if (titleRef.current?.value) {
      postTask(titleRef.current.value);
      setIsAddTaskFlag(false);
    }
  };

  const handleEditTitle = (id: string, title: string) => {
    console.log('Edit title');
    const task = tasks.find((task) => task.id === id);
    if (task) {
      updateTask(id, title, task.isDone);
    }
  };

  const handleEditIsDone = (id: string, isDone: boolean) => {
    const task = tasks.find((task) => task.id === id);
    if (task) {
      updateTask(id, task.title, isDone);
    }
  };

  const handleDelete = (id: string) => {
    deleteTask(id);
  };

  return (
    <TableContainer component={Paper}>
      <Table
        sx={{ minWidth: 650 }}
        aria-label="simple table"
      >
        <TableHead>
          <TableRow>
            <TableCell>Title</TableCell>
            <TableCell align="right">isDone</TableCell>
            <TableCell align="right">createdAt</TableCell>
            <TableCell align="right">delete</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {tasks.map((task) => (
            <TableRow
              key={task.id}
              sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
            >
              <CustomTableCell
                value={task.title}
                handleSave={(value) => handleEditTitle(task.id, value)}
              />
              <CustomTableCellPulldown<boolean>
                value={task.isDone}
                list={[true, false]}
                handleSave={(value) => handleEditIsDone(task.id, value)}
              />
              <TableCell align="right">{task.createdAt}</TableCell>
              <TableCell align="right">
                <IconButton
                  aria-label="delete"
                  onClick={() => handleDelete(task.id)}
                >
                  <Delete />
                </IconButton>
              </TableCell>
            </TableRow>
          ))}

          {isAddTaskFlag ? (
            <TableRow
              key="add-task-row"
              sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
            >
              <TableCell
                component="th"
                scope="row"
              >
                <TextField
                  id="title"
                  label="title"
                  inputRef={titleRef}
                  onBlur={onBlurTitle}
                />
              </TableCell>
              <TableCell align="right"></TableCell>
              <TableCell align="right"></TableCell>
            </TableRow>
          ) : (
            <TableRow
              key="add-task-row"
              sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
            >
              <TableCell
                component="th"
                scope="row"
              >
                <Button
                  variant="contained"
                  color="primary"
                  onClick={() => setIsAddTaskFlag(true)}
                >
                  Add Task
                </Button>
              </TableCell>
              <TableCell align="right"></TableCell>
              <TableCell align="right"></TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default TaskTable;
