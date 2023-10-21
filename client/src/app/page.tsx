'use client';

import { TextField } from '@mui/material';
import { useEffect, useState } from 'react';
import Header from '@/components/Header';
import TaskTable from '@/components/TaskTable';
import customAxios from '@/utils/axios';
import type { Task } from '@/types/tasks';
import type { FC } from 'react';

const HomePage: FC = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [searchText, setSearchText] = useState<string>('');
  const foundTasks = tasks.filter((task) => task.title.includes(searchText));

  const getTasks = async (): Promise<void> => {
    const res = await customAxios.get<Task[]>('/tasks', {
      withCredentials: true,
    });

    setTasks(res.data);
  };

  const postTask = async (title: string): Promise<void> => {
    await customAxios.post<void>(
      '/tasks',
      {
        title: title,
      },
      { withCredentials: true }
    );

    getTasks();
  };

  const updateTask = async (id: string, title: string, isDone: boolean) => {
    await customAxios.put<void>(
      `/tasks/${id}`,
      {
        title: title,
        is_done: isDone,
      },
      { withCredentials: true }
    );

    getTasks();
  };

  const deleteTask = async (id: string): Promise<void> => {
    await customAxios.delete<void>(`/tasks/${id}`, {
      withCredentials: true,
    });

    getTasks();
  };

  const onChangeQuery = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchText(e.target.value);
  };

  useEffect(() => {
    getTasks();
  }, []);

  return (
    <>
      <Header>
        <TextField
          id="query"
          label="query"
          onChange={onChangeQuery}
        />
      </Header>
      <TaskTable
        tasks={foundTasks}
        postTask={postTask}
        updateTask={updateTask}
        deleteTask={deleteTask}
      />
    </>
  );
};

export default HomePage;
