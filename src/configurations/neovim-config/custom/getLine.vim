func! GetSameIndentLine(direction)
  let [inc, endOfFile] = s:getProps(a:direction)

  let lineMarker = line('.')
  let existsSameIndent = 0
  let originalIndent = indent('.')

  while !existsSameIndent && lineMarker != endOfFile
    let lineMarker += inc
    let existsSameIndent = indent(lineMarker) == originalIndent

    if existsSameIndent && IsEmptyLine(lineMarker)
      let existsSameIndent = 0
    endif
  endwhile

  if (existsSameIndent)
    return lineMarker
  else
    return line('.')
  endif
endfunc

func! GetLastMatchingIndent(inc, endOfFile)
  let [inc, endOfFile] = s:getProps(a:direction)

  let lineMarker = line('.')
  let originalIndent = indent('.')

  while lineMarker != endOfFile
    let lineMarker += inc
    let existsSameIndent = indent(lineMarker) == originalIndent

    if existsSameIndent || IsEmptyLine(lineMarker)
      continue
    endif

    break
  endwhile

  return lineMarker
endfunc

func! GetSwitchLine(direction)
  let [inc, endOfFile] = s:getProps(a:direction)

  let lineMarker = line('.')
  let s:isStartEmpty = IsEmptyLine(lineMarker)
  let originalIndent = indent('.')

  func! s:didSwitch(line)
    if s:isStartEmpty == 1
      return !IsEmptyLine(a:line)
    else
      return IsEmptyLine(a:line)
    endif
  endfunc

  while !s:didSwitch(lineMarker) && lineMarker != endOfFile
    let lineMarker += inc
  endwhile

  return lineMarker
endfunc

" [inc, endOfFile]
func! s:getProps(direction)
  if a:direction == 'j'
    return [+1, line('$')]
  elseif a:direction == 'k'
    return [-1, 1]
  endif
endfunc
