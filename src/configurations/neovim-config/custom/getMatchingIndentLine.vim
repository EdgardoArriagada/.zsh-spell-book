func! s:getSameIndentLine(inc, endOfFile)
  let lineMarker = line('.')
  let existsSameIndent = 0
  let originalIndent = indent('.')

  while !existsSameIndent && lineMarker != a:endOfFile
    let lineMarker += a:inc
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

func! s:getLastMatchingIndentLine(inc, endOfFile)
  let lineMarker = line('.')
  let originalIndent = indent('.')

  while lineMarker != a:endOfFile
    let lineMarker += a:inc
    let existsSameIndent = indent(lineMarker) == originalIndent

    if existsSameIndent || IsEmptyLine(lineMarker)
      continue
    endif

    break
  endwhile

  return lineMarker
endfunc

func! GetLastMatchingIndentDown()
  return s:getLastMatchingIndentLine(+1, line('$'))
endfunc!

func! GetLastMatchingIndentUp()
  return s:getLastMatchingIndentLine(-1, 1)
endfunc!

func! GetSameIndentLineDown()
  return s:getSameIndentLine(+1, line('$'))
endfunc!

func! GetSameIndentLineUp()
  return s:getSameIndentLine(-1, 1)
endfunc!

