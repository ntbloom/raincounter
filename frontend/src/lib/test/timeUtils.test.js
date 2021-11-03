const { default: TimeUtils } = require('../data/timeUtils');

function sums(num1, num2) {
  return num1 + num2;
}

describe.each([
  [0, '<1 hour'],
  [55, '<1 hour'],
  [961, '<1 hour'],
  [4501, '1 hour'],
  [7199, '1 hour'],
  [7200, '2 hours'],
  [82888, '23 hours'],
  [86400, '1 day'],
  [172801, '2 days'],
  [701299, '8 days'],
])('test parsing seconds into clean string', (input, exp) => {
  test(`${input} -> ${exp}`, () => {
    expect(TimeUtils.secondsToString(input)).toBe(exp);
  });
});
