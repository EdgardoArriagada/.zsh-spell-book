nnoremap Q :call PowerSelection()<cr>

func! s:hasChar(inputChar)
  return match(getline("."), a:inputChar) > 0
endfunc

func! s:hasBoth(left, right)
  return s:hasChar(a:left) && s:hasChar(a:right)
endfunc

func! s:didSelectInline()
  return col('.') != col('v') && line('.') == line('v')
endfunc

" The order of the elements are the prioritization of the surrounds
let g:pairList = [['(', ')'], ['<', '>'], ['[', ']'], ['{', '}']]

func! PowerSelection()
  let l:savedPos = getpos('.')

  for [l:left, l:right] in g:pairList
    if ! s:hasBoth(l:left, l:right)
      continue
    endif

    execute "normal! \<esc> vi" . l:left

    if s:didSelectInline()
      return
    endif
    
    :call setpos('.', l:savedPos)
  endfor

  for [l:left, l:right] in g:pairList
    if ! s:hasBoth(l:left, l:right)
      continue
    endif

    execute "normal! \<esc>f".l:right."F".l:left."vi".l:left

    if s:didSelectInline()
      return
    endif
    :call setpos('.', l:savedPos)
  endfor

  for [l:left, l:right] in g:pairList
    if ! s:hasBoth(l:left, l:right)
      continue
    endif

    execute "normal! \<esc>0f".l:left."vi".l:left

    if s:didSelectInline()
      return
    endif
    :call setpos('.', l:savedPos) 
  endfor
endfunc

