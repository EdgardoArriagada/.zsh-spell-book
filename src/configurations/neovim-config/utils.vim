func! Strip(input)
  return substitute(a:input, '^\s*\(.\{-}\)\s*$', '\1', '')
endfunc

func! IsEmptyLine(line)
  let lineContent = getline(a:line)
  if lineContent == ''
    return 1
  endif
  return len(Strip(lineContent)) == 0
endfunc

func! IsCursorWithinBuffer()
  return 1 < line('.') && line('.') < line('$')
endfunc

func! JumpUntilNotEmptyLine(direction)
  while IsEmptyLine('.') && IsCursorWithinBuffer()
    execute "normal! ".a:direction."^"
  endwhile
endfunc

