#!/usr/bin/env ruby

require 'pmap'

GIT_MONITOR_PREFIX = File.join Dir.home, '.git-monitor'

def list_respositories
  puts Dir[File.join GIT_MONITOR_PREFIX, '*'].map { |e| File.basename e }
end

def add_repository(url)
  path = url.match(%r{.*/(.*).git})[1]

  Dir.chdir GIT_MONITOR_PREFIX do
    `git clone --depth 1 --quiet --no-single-branch --no-checkout #{url}`
  end

  puts "Cloned repository #{url} into #{path}"
end

def check_repository(path, length: nil)
  length ||= path.length
  full_path = File.join GIT_MONITOR_PREFIX, path
  messages = nil

  update_response = `cd #{full_path} && git remote -v update 2>&1`

  branches_file = File.join GIT_MONITOR_PREFIX, path, '.git-monitor-branches'

  branches = IO.readlines(branches_file).map(&:strip) if File.exist? branches_file
  branches ||= ['master']

  branches.each do |branch|
    next unless update_response.match(/[[:alnum:]]+..[[:alnum:]]+\s+#{branch}/)

    url = `cd #{full_path} && git remote show origin`.match(/Fetch URL: (.*)/) do |m|
      " #{m[1]}"
    end

    messages ||= []
    messages << "#{path.rjust length} on #{branch}: " + 'CHANGED'.rjust(7) + (url || '<unable to get url>')
  end

  messages.join "\n" if messages

  if update_response.match(/error/)
    messages = [
      "#{path.rjust length} ERROR:",
      update_response,
      '-' * 80,
    ].join "\n"
  end

  messages
end

def check_repositories
  repos = Dir[File.join GIT_MONITOR_PREFIX, '*'].map { |e| File.basename e }
  length = repos.map(&:length).max
  changes = repos.pmap { |repo| check_repository repo, length: length }.compact

  if changes.length == 0
    puts 'Already up-to-date.'
  else
    puts changes.join "\n"
  end
end

if ARGV.count == 0
  check_repositories
elsif ARGV.count == 1
  if ARGV.first == '--list' || ARGV.first == '-l'
    list_respositories
  else
    add_repository ARGV.first
  end
else
  puts 'Wrong usage.'
  exit 1
end
