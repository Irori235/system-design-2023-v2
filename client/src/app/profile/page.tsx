'use client';

import { Button, Stack, TextField, Typography } from '@mui/material';
import { useEffect, useRef, useState } from 'react';
import Header from '@/components/Header';
import customAxios from '@/utils/axios';
import type { FC } from 'react';

interface User {
  id: string;
  name: string;
  updatedAt: string;
  createdAt: string;
}

const ProfilePage: FC = () => {
  const refs = useRef<HTMLInputElement[]>([]);
  const [user, setUser] = useState<User | null>(null);

  const patchUpdateName = async (name: string): Promise<void> => {
    const res = await customAxios.patch<void>(
      '/users/name',
      {
        name: name,
      },
      {
        withCredentials: true,
      }
    );

    if (res.status !== 200) {
      console.log(res.statusText);
    }
  };

  const onBlurName = () => {
    if (refs.current[0].value) {
      patchUpdateName(refs.current[0].value);
    }
  };

  const patchUpdatePassword = async (password: string): Promise<void> => {
    try {
      const res = await customAxios.patch<void>('/users/password', {
        password: password,
      });

      if (res.status !== 200) {
        console.log(res.statusText);
      }
    } catch (e: unknown) {
      console.log(e);
    }
  };

  const onBlurPassword = () => {
    if (refs.current[1].value) {
      patchUpdatePassword(refs.current[1].value);
    }
  };

  const getUser = async (): Promise<User> => {
    try {
      const res = await customAxios.get<User>('/users/me', {
        withCredentials: true,
      });

      if (res.status !== 200) {
        console.log(res.statusText);
      }

      return res.data;
    } catch (e: unknown) {
      console.log(e);
      return Promise.reject(new Error('error'));
    }
  };

  const handleLogOut = async (): Promise<void> => {
    try {
      const res = await customAxios.post<void>('/auth/signout', {
        withCredentials: true,
      });

      if (res.status !== 200) {
        console.log(res.statusText);
      }
    } catch (e: unknown) {
      console.log(e);
    }
  };

  const handleDeleteAccount = async (): Promise<void> => {
    try {
      const res = await customAxios.delete<void>('/users/quit', {
        withCredentials: true,
      });

      if (res.status !== 200) {
        console.log(res.statusText);
      }
    } catch (e: unknown) {
      console.log(e);
    }
  };

  const onClickLogOut = () => {
    handleLogOut();
  };

  const onClickDeleteAccount = () => {
    handleDeleteAccount();
  };

  useEffect(() => {
    getUser()
      .then((user) => setUser(user))
      .catch((e) => console.log(e));
  }, []);

  return (
    <>
      <Stack spacing={2}>
        <Header />
        <Typography variant="h1">Profile</Typography>
        <Stack spacing={2}>
          <Typography>Account ID</Typography>
          <Typography>{user?.id}</Typography>
        </Stack>
        <Stack spacing={2}>
          <Typography>Account Name</Typography>
          <TextField
            id="name"
            inputRef={(el) => (refs.current[0] = el)}
            onBlur={onBlurName}
            defaultValue={user?.name}
          />
        </Stack>
        <Stack spacing={2}>
          <Typography>Account Password</Typography>
          <TextField
            id="password"
            inputRef={(el) => (refs.current[1] = el)}
            onBlur={onBlurPassword}
          />
        </Stack>
        <Button
          variant="contained"
          onClick={onClickLogOut}
        >
          LogOut
        </Button>
        <Button
          variant="contained"
          color="primary"
          onClick={onClickDeleteAccount}
        >
          Delete Account
        </Button>
      </Stack>
    </>
  );
};

export default ProfilePage;
