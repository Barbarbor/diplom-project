import React, { useState } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend } from 'recharts';
import { Box, Typography, LinearProgress, Button, List, ListItem } from '@mui/material';
import { QuestionStats } from '@/types/stats';
import { Block } from '../common/Block';
export const SingleChoiceStats = ({ question }: { question: QuestionStats }) => {
  const totalAnswers = question.answers.length;
  const answerCounts = question.answers.reduce((acc, ans) => {
    const numAns = Number(ans);
    acc[numAns] = (acc[numAns] || 0) + 1;
    return acc;
  }, {} as Record<number, number>);

  return (
    <Block>
      <Typography variant="h6" fontWeight="bold" mb={2}>
        {question.label}
      </Typography>
      {question?.options?.map((option) => {
        const count = answerCounts[option.id] || 0;
        const percentage = totalAnswers > 0 ? (count / totalAnswers) * 100 : 0;
        return (
          <Box key={option.id} display="flex" alignItems="center" mb={2}>
            <Typography width="150px" mr={2}>
              {option.label}
            </Typography>
            <Box flex={1}>
              <LinearProgress
                variant="determinate"
                value={percentage}
                sx={{
                  height: 10,
                  borderRadius: 5,
                  backgroundColor: '#e0e0e0',
                  '& .MuiLinearProgress-bar': {
                    backgroundColor: '#1976d2',
                    transition: 'width 0.5s ease-in-out',
                  },
                }}
              />
              <Typography variant="caption" textAlign="center" mt={0.5}>
                {percentage.toFixed(1)}%
              </Typography>
            </Box>
          </Box>
        );
      })}
    </Block>
  );
};export const MultiChoiceStats = ({ question }: { question: QuestionStats }) => {
  const totalInterviews = question.answers.length;
  const answerCounts = question.answers.reduce((acc, ans) => {
    const parsedAns = JSON.parse(ans) as number[];
    parsedAns.forEach((id) => {
      acc[id] = (acc[id] || 0) + 1;
    });
    return acc;
  }, {} as Record<number, number>);

  return (
    <Block>
      <Typography variant="h6" fontWeight="bold" mb={2}>
        {question.label}
      </Typography>
      {question?.options?.map((option) => {
        const count = answerCounts[option.id] || 0;
        const percentage = totalInterviews > 0 ? (count / totalInterviews) * 100 : 0;
        return (
          <Box key={option.id} display="flex" alignItems="center" mb={2}>
            <Typography width="150px" mr={2}>
              {option.label}
            </Typography>
            <Box flex={1}>
              <LinearProgress
                variant="determinate"
                value={percentage}
                sx={{
                  height: 10,
                  borderRadius: 5,
                  backgroundColor: '#e0e0e0',
                  '& .MuiLinearProgress-bar': {
                    backgroundColor: '#2e7d32',
                    transition: 'width 0.5s ease-in-out',
                  },
                }}
              />
              <Typography variant="caption" textAlign="center" mt={0.5}>
                {percentage.toFixed(1)}%
              </Typography>
            </Box>
          </Box>
        );
      })}
    </Block>
  );
};export const ConsentStats = ({ question }: { question: QuestionStats }) => {
  const totalAnswers = question.answers.length;
  const trueCount = question.answers.filter((ans) => ans === 'true').length;
  const falseCount = totalAnswers - trueCount;
  const truePercentage = totalAnswers > 0 ? (trueCount / totalAnswers) * 100 : 0;
  const falsePercentage = totalAnswers > 0 ? (falseCount / totalAnswers) * 100 : 0;

  return (
    <Block>
      <Typography variant="h6" fontWeight="bold" mb={2}>
        {question.label}
      </Typography>
      <Box display="flex" alignItems="center" mb={2}>
        <Typography width="150px" mr={2}>
          Согласны
        </Typography>
        <Box flex={1}>
          <LinearProgress
            variant="determinate"
            value={truePercentage}
            sx={{
              height: 10,
              borderRadius: 5,
              backgroundColor: '#e0e0e0',
              '& .MuiLinearProgress-bar': {
                backgroundColor: '#ab47bc',
                transition: 'width 0.5s ease-in-out',
              },
            }}
          />
          <Typography variant="caption" textAlign="center" mt={0.5}>
            {truePercentage.toFixed(1)}%
          </Typography>
        </Box>
      </Box>
      <Box display="flex" alignItems="center" mb={2}>
        <Typography width="150px" mr={2}>
          Не согласны
        </Typography>
        <Box flex={1}>
          <LinearProgress
            variant="determinate"
            value={falsePercentage}
            sx={{
              height: 10,
              borderRadius: 5,
              backgroundColor: '#e0e0e0',
              '& .MuiLinearProgress-bar': {
                backgroundColor: '#d32f2f',
                transition: 'width 0.5s ease-in-out',
              },
            }}
          />
          <Typography variant="caption" textAlign="center" mt={0.5}>
            {falsePercentage.toFixed(1)}%
          </Typography>
        </Box>
      </Box>
    </Block>
  );
};
export const TextStats = ({ question }: { question: QuestionStats }) => {
  const [showAll, setShowAll] = useState(false);

  const answerCounts = question.answers.reduce((acc, ans) => {
    acc[ans] = (acc[ans] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  const sortedAnswers = Object.entries(answerCounts).sort((a, b) => b[1] - a[1]);
  const top5 = sortedAnswers.slice(0, 5);
  const remaining = sortedAnswers.slice(5);

  return (
    <Block>
      <Typography variant="h6" fontWeight="bold" mb={2}>
        {question.label}
      </Typography>
      <List>
        {top5.map(([answer, count]) => (
          <ListItem key={answer} sx={{ py: 1 }}>
            <Typography variant="body1">{answer} ({count} раз)</Typography>
          </ListItem>
        ))}
      </List>
      {remaining.length > 0 && (
        <Button
          variant="text"
          color="primary"
          onClick={() => setShowAll(!showAll)}
          sx={{ mt: 2 }}
        >
          {showAll ? 'Скрыть' : 'Показать остальные ответы'}
        </Button>
      )}
      {showAll && (
        <List sx={{ mt: 2 }}>
          {remaining.map(([answer, count]) => (
            <ListItem key={answer} sx={{ py: 1 }}>
              <Typography variant="body1">{answer} ({count} раз)</Typography>
            </ListItem>
          ))}
        </List>
      )}
    </Block>
  );
};export const DateStats = ({ question }: { question: QuestionStats }) => {
  const dates = question.answers.map((ans) => new Date(ans)).filter((date) => !isNaN(date.getTime()));
  const frequencyMap = dates.reduce((acc, date) => {
    const dateStr = date.toISOString().split('T')[0];
    acc[dateStr] = (acc[dateStr] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  const data = Object.entries(frequencyMap)
    .map(([date, count]) => ({ date, count }))
    .sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());

  const maxCount = Math.max(...data.map((d) => d.count), 1);
  const tickCount = maxCount;

  return (
    <Block>
      <h3 className="font-semibold text-lg mb-2">{question.label}</h3>
      <LineChart
        width={600}
        height={300}
        data={data}
        margin={{ top: 40, right: 30, left: 20, bottom: 30 }}
      >
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis
          dataKey="date"
          label={{ value: 'Даты', position: 'insideBottomRight', offset: -10 }}
          padding={{ left: 20 }}
        />
        <YAxis
          domain={[1, maxCount]}
          tickCount={tickCount}
          interval="preserveStartEnd"
          tickFormatter={(value) => Math.round(value).toString()}
          label={{
            value: 'Количество ответов',
            angle: -90,
            position: 'insideLeft',
            offset: -10,
            style: { textAnchor: 'middle' },
          }}
        />
        <Tooltip />
        <Legend />
        <Line type="monotone" dataKey="count" stroke="#82ca9d" />
      </LineChart>
    </Block>
  );
};
export const NumberStats = ({ question }: { question: QuestionStats }) => {
  const numbers = question.answers.map(Number).filter((num) => !isNaN(num) && num !== 0);
  const frequencyMap = numbers.reduce((acc, num) => {
    acc[num] = (acc[num] || 0) + 1;
    return acc;
  }, {} as Record<number, number>);

  const data = Object.entries(frequencyMap)
    .map(([num, count]) => ({ number: Number(num), count }))
    .sort((a, b) => a.number - b.number);

  const maxCount = Math.max(...data.map((d) => d.count), 1);
  const tickCount = maxCount;

  return (
    <Block>
      <h3 className="font-semibold text-lg mb-2">{question.label}</h3>
      <LineChart
        width={600}
        height={300}
        data={data}
        margin={{ top: 40, right: 30, left: 20, bottom: 30 }}
      >
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis
          dataKey="number"
          domain={[1, 'dataMax']}
          label={{ value: 'Числа', position: 'insideBottomRight', offset: -10 }}
          padding={{ left: 20 }}
        />
        <YAxis
          domain={[1, maxCount]}
          tickCount={tickCount}
          interval="preserveStartEnd"
          tickFormatter={(value) => Math.round(value).toString()}
          label={{
            value: 'Количество ответов',
            angle: -90,
            position: 'insideLeft',
            offset: -10,
            style: { textAnchor: 'middle' },
          }}
        />
        <Tooltip />
        <Legend />
        <Line type="monotone" dataKey="count" stroke="#8884d8" />
      </LineChart>
    </Block>
  );
};
export const RatingStats = ({ question }: { question: QuestionStats }) => {
  const totalAnswers = question.answers.length;
  const sum = question.answers.reduce((acc, ans) => acc + Number(ans), 0);
  const average = totalAnswers > 0 ? (sum / totalAnswers).toFixed(1) : '0.0';

  const frequency = question.answers.reduce((acc, ans) => {
    const numAns = Number(ans);
    acc[numAns] = (acc[numAns] || 0) + 1;
    return acc;
  }, {} as Record<number, number>);

  const maxRating = 5;

  const data = Array.from({ length: maxRating }, (_, i) => i + 1).map((star) => ({
    rating: star,
    count: frequency[star] || 0,
  }));

  const maxCount = Math.max(...data.map((d) => d.count), 1);
  const tickCount = maxCount + 1;

  return (
    <Block>
      <h3 className="font-semibold text-lg mb-2">{question.label}</h3>
      <p className="mb-4">Средний рейтинг: {average}</p>
      <LineChart
        width={600}
        height={300}
        data={data}
        margin={{ top: 40, right: 30, left: 20, bottom: 30 }}
      >
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis
          dataKey="rating"
          domain={[1, maxRating]}
          label={{ value: 'Рейтинг (звёзды)', position: 'insideBottomRight', offset: -10 }}
          padding={{ left: 20 }}
        />
        <YAxis
          domain={[0, maxCount]}
          tickCount={tickCount}
          interval="preserveStartEnd"
          tickFormatter={(value) => Math.round(value).toString()}
          label={{
            value: 'Количество ответов',
            angle: -90,
            position: 'insideLeft',
            offset: -10,
            style: { textAnchor: 'middle' },
          }}
        />
        <Tooltip />
        <Legend />
        <Line type="monotone" dataKey="count" stroke="#ffbb28" />
      </LineChart>
    </Block>
  );
};
export const EmailStats = ({ question }: { question: QuestionStats }) => {
  const [showAll, setShowAll] = useState(false);

  const uniqueEmails = [...new Set(question.answers)].sort();
  const top5 = uniqueEmails.slice(0, 5);
  const remaining = uniqueEmails.slice(5);

  return (
    <Block>
      <Typography variant="h6" fontWeight="bold" mb={2}>
        {question.label}
      </Typography>
      <List>
        {top5.map((email) => (
          <ListItem key={email} sx={{ py: 1 }}>
            <Typography variant="body1">{email}</Typography>
          </ListItem>
        ))}
      </List>
      {remaining.length > 0 && (
        <Button
          variant="text"
          color="primary"
          onClick={() => setShowAll(!showAll)}
          sx={{ mt: 2 }}
        >
          {showAll ? 'Скрыть' : 'Показать остальные email'}
        </Button>
      )}
      {showAll && (
        <List sx={{ mt: 2 }}>
          {remaining.map((email) => (
            <ListItem key={email} sx={{ py: 1 }}>
              <Typography variant="body1">{email}</Typography>
            </ListItem>
          ))}
        </List>
      )}
    </Block>
  );
};