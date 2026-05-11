export const getArchitects = () => {
  // Get all <a> tags (architects)
  const architects = Array.from(document.getElementsByTagName('a'))
  
  // Get all elements except <a> tags (non-architects)
  const nonArchitects = Array.from(document.querySelectorAll('*:not(a)'))
    .filter(el => el.parentElement && el.parentElement !== document.documentElement)
  
  return [architects, nonArchitects]
}

export const getClassical = () => {
  // Get all architects first
  const [architects] = getArchitects()
  
  // Filter architects with 'classical' class
  const classicalArchitects = architects.filter(el => el.classList.contains('classical'))
  
  // Non-classical architects
  const nonClassicalArchitects = architects.filter(el => !el.classList.contains('classical'))
  
  return [classicalArchitects, nonClassicalArchitects]
}

export const getActive = () => {
  // Get all classical architects
  const [classicalArchitects] = getClassical()
  
  // Filter for active classical architects
  const activeClassical = classicalArchitects.filter(el => el.classList.contains('active'))
  
  // Non-active classical architects
  const nonActiveClassical = classicalArchitects.filter(el => !el.classList.contains('active'))
  
  return [activeClassical, nonActiveClassical]
}

export const getBonannoPisano = () => {
  // Find the element with id 'BonannoPisano'
  const bonannoPisano = document.getElementById('BonannoPisano')
  
  // Get all active classical architects
  const [activeClassical] = getActive()
  
  // Filter out BonannoPisano from active classical architects
  const remaining = activeClassical.filter(el => el.id !== 'BonannoPisano')
  
  return [bonannoPisano, remaining]
}
