import { TextField } from '@mui/material';
import { useEffect, useState } from 'react';
import customAxios from '@/utils/axios';
import type { Task } from '@/types/tasks';
import type { AxiosResponse } from 'axios';
import type { FC } from 'react';

const SearchBox: FC = () => {
  const [isSearchTaskActive, setIsSearchTaskActive] = useState<boolean>(false);
  const [query, setQuery] = useState<string>('');
  const [searchedTasks, setSearchedTasks] = useState<Task[]>([]);

  const getSearchedTasks = async (query: string): Promise<void> => {
    try {
      const res: AxiosResponse<Task[]> = await customAxios.get<Task[]>(
        `/search?q=${query}`,
        { withCredentials: true }
      );
      setSearchedTasks(res.data);
    } catch (e: unknown) {
      console.log(e);
    }
  };

  useEffect(() => {
    if (isSearchTaskActive) {
      getSearchedTasks(query);
    }
  }, [isSearchTaskActive, query]);

  const onChangeQuery = (e: React.ChangeEvent<HTMLInputElement>) => {
    setQuery(e.target.value);
  };

  return (
    <TextField
      id="query"
      label="query"
      onChange={onChangeQuery}
    />
  );
};

export default SearchBox;
