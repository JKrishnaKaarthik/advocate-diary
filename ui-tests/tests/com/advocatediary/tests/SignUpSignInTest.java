package com.advocatediary.tests;

import com.advocatediary.pages.DiaryPage;
import com.advocatediary.pages.SignInPage;
import com.advocatediary.pages.SignUpPage;
import org.testng.Assert;
import org.testng.annotations.Test;

public class SignUpSignInTest extends BaseTest {

  @Test
  public void existingLawyerSignsInFromWelcomeBack() {
    SignInPage signIn = new SignInPage(driver).open(baseUrl);
    Assert.assertEquals(signIn.heading(), "Welcome back");

    DiaryPage diary = signIn.signIn(LAWYER_EMAIL, LAWYER_PASSWORD);
    Assert.assertEquals(diary.heading(), "Court diary");
  }

  @Test
  public void shortPasswordShowsHintOnSignUp() {
    SignUpPage signUp = new SignUpPage(driver).open(baseUrl);
    signUp.submitInvalid("Test Advocate", "short@chamber.test", "short");
    Assert.assertTrue(signUp.hintText().toLowerCase().contains("at least 8"));
  }
}
