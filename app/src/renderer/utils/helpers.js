export const formatDateString = (dateString, getyear, getmonth, getday) => {
  // Create a Date object
  const dateObject = new Date(dateString)

  // Extract different parts of the date
  const day = dateObject.getDate() // Get the day of the month
  const month = dateObject.getMonth() + 1 // Get the month (0-11), so add 1
  const year = dateObject.getFullYear() // Get the year

  // Format the date as needed, for example in YYYY-MM-DD format
  if (getyear && !getmonth && !getday) return year

  if (getyear && getmonth && !getday)
    return `${year}-${month < 10 ? '0' + month : month}`

  if (getyear && getmonth && getday)
    return `${year}-${month < 10 ? '0' + month : month}-${day < 10 ? '0' + day : day}`
}

export const toMinutes = (number) => {
  number = number * 60
  const minutes = Math.floor(number / 60) // 7
  let seconds = ((number % 60) / 100) * 60 // 30
  seconds = Math.round(seconds)
  return minutes + ' m, ' + seconds + ' s'
}

export default formatDateString
