# brew install ffmpeg
# brew install gifsicle

convertMovToGif() {
  local input=${1:?'input file is required'}
  local output=${input:r}.gif

  if [[ ! -f ${input} ]]
    then ${zsb}.throw "File ${input} does not exist"
  fi

  ffmpeg -i ${input} -pix_fmt rgb8 -r 10 ${output} && gifsicle -O3 ${output} -o ${output}
}

_${zsb}.convertMovToGif() {
  local filesWithMovExtension=( ${(@f)$(ls | rg '\.mov$')} )

  _describe 'command' filesWithMovExtension
}

compdef _${zsb}.convertMovToGif convertMovToGif
