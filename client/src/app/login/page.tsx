'use client';
import { Stack, Typography, TextField, Button } from '@mui/material';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import customAxios from '@/utils/axios';
import type { AxiosResponse } from 'axios';
import type { FC } from 'react';

const SignInPage: FC = () => {
  const [name, setname] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const router = useRouter();

  const onChangename = (e: React.ChangeEvent<HTMLInputElement>) => {
    setname(e.target.value);
  };

  const onChangePassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
  };

  const postSignIn = async (name: string, password: string): Promise<void> => {
    try {
      const res: AxiosResponse<void> = await customAxios.post('/auth/signin', {
        name: name,
        password: password,
      });

      if (res.status === 200) {
        router.push('/');
      }
    } catch (e: unknown) {
      console.log(e);
    }
  };

  const handleSignIn = () => {
    postSignIn(name, password);
  };

  return (
    <Stack spacing={2}>
      <Typography variant="h1">LogIn</Typography>
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
        onClick={handleSignIn}
      >
        SignIn
      </Button>
      <Link href="/signup">SignUp</Link>
    </Stack>
  );
};

export default SignInPage;
