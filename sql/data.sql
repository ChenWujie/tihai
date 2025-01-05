CREATE TABLE `users` (
     `id` INT AUTO_INCREMENT PRIMARY KEY,
     `username` VARCHAR(50) NOT NULL UNIQUE,
     `password` VARCHAR(255) NOT NULL,
     `role` ENUM('student', 'teacher', 'admin') NOT NULL,
     `email` VARCHAR(100) DEFAULT NULL,
     `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
     `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     `deleted_at` DATETIME DEFAULT NULL
);

CREATE TABLE `questions` (
     `id` INT AUTO_INCREMENT PRIMARY KEY,
     `title` TEXT NOT NULL,
     `content` TEXT,
     `image_url` VARCHAR(255),
     `teacher_id` INT NOT NULL,
     `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
     `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     `deleted_at` DATETIME DEFAULT NULL,
     FOREIGN KEY (`teacher_id`) REFERENCES `users` (`id`)
);

CREATE TABLE `student_answers` (
     `id` INT AUTO_INCREMENT PRIMARY KEY,
     `student_id` INT NOT NULL,
     `question_id` INT NOT NULL,
     `answer_text` TEXT,
     `answer_image_url` VARCHAR(255),
     `submit_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
     FOREIGN KEY (`student_id`) REFERENCES `users` (`id`),
     FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`)
);

CREATE TABLE `scores` (
      `id` INT AUTO_INCREMENT PRIMARY KEY,
      `student_id` INT NOT NULL,
      `question_id` INT NOT NULL,
      `score` INT,
      `graded_by` INT,
      `graded_time` DATETIME,
      FOREIGN KEY (`student_id`) REFERENCES `users` (`id`),
      FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`),
      FOREIGN KEY (`graded_by`) REFERENCES `users` (`id`)
);

CREATE TABLE `comments` (
      `id` INT AUTO_INCREMENT PRIMARY KEY,
      `answer_id` INT NOT NULL,
      `teacher_id` INT NOT NULL,
      `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
      `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      `deleted_at` DATETIME DEFAULT NULL,
      FOREIGN KEY (`answer_id`) REFERENCES `student_answers` (`id`),
      FOREIGN KEY (`teacher_id`) REFERENCES `users` (`id`)
);