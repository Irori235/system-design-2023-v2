'use client';

import { Stack, Typography, TextField, Button } from '@mui/material';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import customAxios from '@/utils/axios';
import type { FC } from 'react';

interface ID {
  id: string;
}

const SignUpPage: FC = () => {
  const router = useRouter();
  const [name, setname] = useState<string>('');
  const [password, setPassword] = useState<string>('');

  const onChangename = (e: React.ChangeEvent<HTMLInputElement>) => {
    setname(e.target.value);
  };

  const onChangePassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
  };

  const postSignUp = async (name: string, password: string): Promise<void> => {
    try {
      const res = await customAxios.post<ID>('/auth/signup', {
        name: name,
        password: password,
      });

      if (res.status !== 200) {
        console.log(res.statusText);
      }

      router.push(`/login`);
    } catch (e: unknown) {
      console.log(e);
    }
  };

  const handleSignUp = () => {
    postSignUp(name, password);
  };

  return (
    <Stack spacing={2}>
      <Typography variant="h1">SignUp</Typography>
      <TextField
        id="name"
        label="name"
        onChange={onChangename}
        required
      />
      <TextField
        id="password"
        label="Password"
        onChange={onChangePassword}
        required
      />
      <Button
        variant="contained"
        color="primary"
        onClick={handleSignUp}
      >
        SignUp
      </Button>
    </Stack>
  );
};

export default SignUpPage;
