import { AccountCircle } from '@mui/icons-material';
import { Box, IconButton, Stack, Typography } from '@mui/material';
import { useRouter } from 'next/navigation';
import type { FC } from 'react';

interface Props {
  children?: React.ReactNode;
}

const Header: FC<Props> = ({ children }) => {
  const router = useRouter();

  return (
    <>
      <Stack
        spacing={2}
        direction="row"
      >
        <Box onClick={() => router.push('/')}>
          <Typography variant="h1">Task Management</Typography>
        </Box>
        {children}
        <IconButton
          aria-label="accountCircle"
          size="large"
          onClick={() => router.push('/profile')}
        >
          <AccountCircle />
        </IconButton>
      </Stack>
    </>
  );
};

export default Header;
