import { greetUser } from '../src/utils/greeting';

describe('greetUser', () => {
  it('should return a greeting message', () => {
    expect(greetUser('John')).toBe('Hello, John!');
  });
});