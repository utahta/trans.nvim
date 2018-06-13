import neovim
import os
import subprocess

@neovim.plugin
class Trans(object):
    def __init__(self, vim):
        self._vim = vim

    @neovim.command('Trans', nargs='?', range='%')
    def trans(self, *args):
        text = self._get_text()
        to = self._opt_trans_lang_locale()
        self._translate(text, to)

    def _translate(self, text, to):
        if text == '':
            return

        cmd = ['trans', '-t', to]
        creds = self._opt_trans_lang_credentials_file()
        if creds != '':
            cmd.extend(['-c', creds])
        cmd.append(text)
        b = subprocess.check_output(cmd)

        output = b.decode('utf-8').strip()
        if self._opt_trans_lang_output() == 'preview':
            self._vim.command("silent pclose")
            self._vim.command("silent pedit +set\\ noswapfile\\ buftype=nofile translated")
            self._vim.command("wincmd P")
            buf = self._vim.current.buffer
            buf[:] = [output]
            self._vim.command("wincmd p")
        else:
            self._vim.out_write("%s\n" % output)

    def _get_text(self):
        start_pos = self._vim.eval("getpos(\"'<\")")
        end_pos = self._vim.eval("getpos(\"'>\")")

        if start_pos[1] == 0 and start_pos[2] == 0 and end_pos[1] == 0 and end_pos[2] == 0:
            return ""

        lines = self._vim.current.buffer[start_pos[1]-1:end_pos[1]]
        for i, line in enumerate(lines):
            if i == 0:
                line = line[start_pos[2]-1:]
            elif i == len(lines)-1:
                if end_pos[2] > len(line):
                    end_pos[2] = len(line)
                line = line[:end_pos[2]]

            line = line.strip()
            for c in self._opt_trans_lang_cutset():
                if line.startswith(c):
                    line = line[len(c):].strip()
            lines[i] = line

        return " ".join(lines)

    def _opt_trans_lang_locale(self):
        if 'trans_lang_locale' in self._vim.vars:
            return self._vim.vars['trans_lang_locale']
        return 'ja'

    def _opt_trans_lang_cutset(self):
        if 'trans_lang_cutset' in self._vim.vars:
            return self._vim.vars['trans_lang_cutset'].split(' ')
        return ['//', '#']

    def _opt_trans_lang_credentials_file(self):
        if 'trans_lang_credentials_file' in self._vim.vars:
            v = self._vim.vars['trans_lang_credentials_file']
            if v.startswith("~"):
                v = v[1:]
                v = '%s/%s' % (os.environ['HOME'], v)
            return v
        return ''

    def _opt_trans_lang_output(self):
        if 'trans_lang_output' in self._vim.vars:
            return self._vim.vars['trans_lang_output']
        return ''
